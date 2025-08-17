package components


// other components to be added later

/*
  templ SortableScript(id string) {
	<script type="module">
    import Sortable from {{ BaseUrl + "/assets/js/sortable.js" }}
    new Sortable(document.getElementById('{{ id }}'), {
        animation: 150,
		handle: '.handle',
        ghostClass: 'opacity-25',
        onEnd: (evt) => {
            document.getElementById('{{ id }}').dispatchEvent(
                new CustomEvent('sorted', {detail: {
					oldIndex: evt.oldIndex,
					newIndex: evt.newIndex
					}})
            )
        }
    })
</script>
}
*/
