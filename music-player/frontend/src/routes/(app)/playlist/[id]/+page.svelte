<h1>Playlist Editor</h1>

<script lang="ts">
  import { page } from "$app/state";
  import { apiFetch } from "$lib/api";
  import PlaylistOverview from "$lib/components/PlaylistOverview.svelte";
  import type { Song, Playlist } from "$lib/types";
  import { onMount } from "svelte";

  const { id } = page.params

  let playlist = $state<Playlist | null>(null);
  let songs = $state<(Song | null)[]>([]);
  let addSongID = $state<number>(1);

  async function fetchPlaylist() {
    try {
      let data: Playlist = await apiFetch<Playlist>("/playlists/" + id);
      playlist = data;
    } catch (error) {
      console.error(error);
    }
  }

  async function fetchSongs() {
    if (songs == null) return
    try {
      let data: Song[] = await apiFetch<Song[]>("/playlists/"+ id +"/songs");
      songs = data;
    } catch (error) {
      console.error(error);
    }
  }

  async function addSong() {
    try {
      const data = await apiFetch("/playlists/"+ id +"/songs/"+ addSongID, {
        method: 'POST',
      })
      console.log(data);
      fetchPlaylist();
    } catch (error) {
      console.error(error)
    }
  }

  async function removeSong(songId: number) {
    try {
      const data = await apiFetch("/playlists/"+ id +"/songs/"+ songId, {
        method: 'DELETE',
      })
      console.log(data);
      fetchPlaylist();
    } catch (error) {
      console.error(error)
    }
  }

  onMount(fetchPlaylist);
</script>

{#if playlist}
  <input type="number" bind:value={addSongID} />
  <button onclick={addSong}>Add Song</button>

  <div style="padding-top: 8px;">
    <PlaylistOverview {playlist} />
  </div>
{:else}
  <div>Couldn't find playlist with ID: "{id}"</div>
{/if}