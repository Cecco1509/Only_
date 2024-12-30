<script lang="ts">
	import { setContext } from 'svelte';
	import '../../app.css';
	let { children, data} = $props();

	setContext('user', data.user)
</script>

{@render children()}