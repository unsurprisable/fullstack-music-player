<script lang="ts">
  import type { Playlist } from "$lib/types";

  export let playlist: Playlist;
  export let editAction: () => void;
  export let deleteAction: () => void;
</script>

<p>======================================================</p>
<h3>{playlist.name}</h3>
<div>{playlist.songIDs.length} songs</div>
<div>created on: {playlist.createdAt}</div>

<div style="padding-top: 10px;">
  <button onclick={editAction}>Edit</button>
  <button onclick={deleteAction}>Delete</button>
</div>

