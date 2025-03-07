<script lang="ts">
  import type { Song } from "$lib/types";
  
  let { song, onclick } = $props<{song: Song, onclick: Function}>();
</script>

<p>======================================================</p>
<h1>{song.title}</h1>
<h3>{song.artist} - {song.album}</h3>
<p>uploaded at: {song.uploadedAt}</p>
<audio controls>
  <source src={song.fileURL} type ="audio/mpeg" />
  Your browser does not support the audio element.
</audio>
<div>
  <button {onclick}>Delete</button>
</div>