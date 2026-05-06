/**
 * <chimp-dropzone> — Web component for drag-and-drop file uploads.
 *
 * Handles file transfer from drop events to a hidden <input type="file">,
 * validates files (type, size, count), renders a file list, and dispatches
 * a "fileschanged" CustomEvent for Datastar signal updates.
 *
 * Uses event delegation on the component element itself to avoid timing
 * issues with Datastar's DOM processing.
 *
 * Data attributes (set by the templ component):
 *   data-max-size   — max bytes per file (0 = no limit)
 *   data-max-files  — max number of files (0 = unlimited)
 *   data-accept     — accepted file types (e.g. "image/*", ".pdf,.doc")
 */
class ChimpDropzone extends HTMLElement {
  connectedCallback() {
    this._dataTransfer = new DataTransfer();

    this.addEventListener('click', (e) => {
      if (e.target.closest('[data-remove]')) return;
      if (e.target.closest('button[type="submit"]')) return;
      const area = this.querySelector('[data-dropzone-area]');
      if (area && (e.target === area || area.contains(e.target))) {
        const input = this.querySelector('input[type="file"]');
        if (input && !input.disabled) input.click();
      }
    });

    this.addEventListener('drop', (e) => {
      const area = this.querySelector('[data-dropzone-area]');
      if (area && (e.target === area || area.contains(e.target))) {
        this._onDrop(e);
      }
    });

    this.addEventListener('change', (e) => {
      if (e.target.matches('input[type="file"]')) {
        this._onInputChange();
      }
    });

    this.addEventListener('click', (e) => {
      const removeBtn = e.target.closest('[data-remove]');
      if (removeBtn) {
        e.stopPropagation();
        this._removeFile(parseInt(removeBtn.dataset.remove));
      }
    });
  }

  get _input() { return this.querySelector('input[type="file"]'); }
  get _filelist() { return this.querySelector('.dropzone-filelist'); }
  get _maxSize() { return parseInt(this.dataset.maxSize) || 0; }
  get _maxFiles() { return parseInt(this.dataset.maxFiles) || 0; }
  get _accept() { return this.dataset.accept || ''; }
  get _multiple() { return this._input?.hasAttribute('multiple'); }

  _onDrop(e) {
    const files = Array.from(e.dataTransfer?.files || []);
    this._processFiles(files);
  }

  _onInputChange() {
    const input = this._input;
    if (!input) return;
    const files = Array.from(input.files);
    this._processFiles(files);
  }

  _processFiles(newFiles) {
    const validated = newFiles.filter(f => this._validateFile(f));

    if (!this._multiple) {
      this._dataTransfer = new DataTransfer();
      if (validated.length > 0) {
        this._dataTransfer.items.add(validated[0]);
      }
    } else {
      for (const f of validated) {
        if (this._maxFiles > 0 && this._dataTransfer.files.length >= this._maxFiles) break;
        this._dataTransfer.items.add(f);
      }
    }

    const input = this._input;
    if (input) input.files = this._dataTransfer.files;
    this._renderFileList();
    this._dispatchChange();
  }

  _validateFile(file) {
    if (this._maxSize > 0 && file.size > this._maxSize) return false;
    if (this._accept && !this._matchAccept(file)) return false;
    return true;
  }

  _matchAccept(file) {
    const accept = this._accept;
    const parts = accept.split(',').map(s => s.trim().toLowerCase());
    const fileName = file.name.toLowerCase();
    const fileType = file.type.toLowerCase();

    return parts.some(part => {
      if (part.startsWith('.')) {
        return fileName.endsWith(part);
      }
      if (part.endsWith('/*')) {
        return fileType.startsWith(part.slice(0, -1));
      }
      return fileType === part;
    });
  }

  _removeFile(index) {
    const dt = new DataTransfer();
    for (let i = 0; i < this._dataTransfer.files.length; i++) {
      if (i !== index) dt.items.add(this._dataTransfer.files[i]);
    }
    this._dataTransfer = dt;
    const input = this._input;
    if (input) input.files = this._dataTransfer.files;
    this._renderFileList();
    this._dispatchChange();
  }

  _renderFileList() {
    const filelist = this._filelist;
    if (!filelist) return;
    const files = this._dataTransfer.files;
    if (files.length === 0) {
      filelist.innerHTML = '';
      return;
    }

    filelist.innerHTML = Array.from(files).map((f, i) =>
      `<div class="flex items-center justify-between gap-2 px-3 py-1.5 text-sm rounded-md bg-muted/50">
        <span class="truncate">${this._escapeHtml(f.name)}</span>
        <span class="flex items-center gap-2 shrink-0">
          <span class="text-muted-foreground text-xs">${this._formatSize(f.size)}</span>
          <button type="button" data-remove="${i}" class="text-muted-foreground hover:text-foreground text-xs cursor-pointer">&times;</button>
        </span>
      </div>`
    ).join('');
  }

  _dispatchChange() {
    this.dispatchEvent(new CustomEvent('fileschanged', {
      bubbles: true,
      detail: { count: this._dataTransfer.files.length }
    }));
  }

  _formatSize(bytes) {
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
  }

  _escapeHtml(str) {
    const div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
  }

  clear() {
    this._dataTransfer = new DataTransfer();
    const input = this._input;
    if (input) input.files = this._dataTransfer.files;
    this._renderFileList();
    this._dispatchChange();
  }
}

customElements.define('chimp-dropzone', ChimpDropzone);
