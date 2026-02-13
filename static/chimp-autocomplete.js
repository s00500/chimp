/**
 * <chimp-autocomplete> — Web component wrapping BaseCoat's select/combobox.
 *
 * Watches the listbox for child mutations (e.g. Datastar patching new options)
 * and automatically reinitializes BaseCoat while preserving the popover state.
 *
 * Usage:
 *   <chimp-autocomplete>
 *     <div class="select w-full">
 *       <button type="button" class="btn-outline ...">...</button>
 *       <div data-popover aria-hidden="true">
 *         <header>...</header>
 *         <div role="listbox">...</div>
 *       </div>
 *       <input type="hidden" name="..." />
 *     </div>
 *   </chimp-autocomplete>
 */
class ChimpAutocomplete extends HTMLElement {
  connectedCallback() {
    this._select = this.querySelector('.select');
    this._ensureInit();

    const listbox = this.querySelector('[role="listbox"]');
    if (listbox) {
      this._observer = new MutationObserver(() => this._onOptionsChanged());
      this._observer.observe(listbox, { childList: true });
    }
  }

  disconnectedCallback() {
    this._observer?.disconnect();
    this._observer = null;
    this._select = null;
  }

  _ensureInit() {
    if (this._select && !this._select.hasAttribute('data-select-initialized')) {
      window.basecoat?.initAll();
    }
  }

  _onOptionsChanged() {
    const popover = this.querySelector('[data-popover]');
    const input = this.querySelector('input[role="combobox"]');
    const wasOpen = popover?.getAttribute('aria-hidden') === 'false';
    const searchVal = input?.value || '';

    if (this._select) {
      this._select.removeAttribute('data-select-initialized');
      window.basecoat?.initAll();
    }

    if (wasOpen) {
      requestAnimationFrame(() => {
        if (popover) popover.setAttribute('aria-hidden', 'false');
        const trigger = this.querySelector('button[aria-haspopup]');
        if (trigger) trigger.setAttribute('aria-expanded', 'true');
        if (input) {
          input.value = searchVal;
          input.focus();
        }
      });
    }
  }
}

customElements.define('chimp-autocomplete', ChimpAutocomplete);
