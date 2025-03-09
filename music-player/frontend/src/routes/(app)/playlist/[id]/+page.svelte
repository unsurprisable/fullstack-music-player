<h1>Playlist Editor</h1>

<script lang="ts">
  import { page } from "$app/state";
  import { apiFetch } from "$lib/api";
  import type { Song, Playlist } from "$lib/types";
  import { onMount } from "svelte";

  const { id } = page.params

  let playlist = $state<Playlist | null>(null);
  let songs = $state<(Song | null)[]>([]);
  let addSongID = $state<number>(1);

  async function fetchPlaylist() {
    try {
      let data: Playlist = await apiFetch<Playlist>("/playlists/" + id);
      console.log(data);
      playlist = data;
      songs = Array(playlist.songIDs.length).fill(null)
      playlist.songIDs.forEach((id, i) => {fetchSong(i, id)})
    } catch (error) {
      console.error(error);
    }
  }

  async function fetchSong(i: number, id: number) {
    if (songs == null) return
    try {
      let data: Song = await apiFetch<Song>("/songs/" + id);
      songs[i] = data;
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
  <h3>{playlist.name}</h3>
  <input type="number" bind:value={addSongID} />
  <button onclick={addSong}>Add Song</button>

  <div>
    {#if songs && songs.length > 0}
      {#each songs as song, i}
        <div style="padding-top: 7px;">
          {i+1}. 
          {#if song}
            {song.artist} - {song.title} [{song.album}]
            <div>
              <button onclick={() => {removeSong(song.id)}}>Remove</button>
            </div>
          {:else}
            <em>Loading song...</em>
          {/if}
        </div>
      {/each}
    {/if}
  </div>
{:else}
  <div>Couldn't find playlist with ID: "{id}"</div>
{/if}