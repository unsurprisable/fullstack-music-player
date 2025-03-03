<script lang="ts">
  import type { Song } from "$lib/types";

  let { song, onclick } = $props<{song: Song, onclick: Function}>(); 
</script>

<button {onclick}>
  {song.artist} {song.artist != "" ? "-" : ""} {song.title} {song.album != "" ? "(" + song.album + ")" : ""}
</button>