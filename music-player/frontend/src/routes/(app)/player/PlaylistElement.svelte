<script lang="ts">
  import type { Playlist } from "$lib/types";

  // let { playlist, onclick } = $props<{song: Playlist, onclick: Function}>();
  export let playlist: Playlist;
  export let onclick: () => void;
</script>

<button {onclick}>
  {playlist.name} - {playlist.songIDs.length} songs
</button>