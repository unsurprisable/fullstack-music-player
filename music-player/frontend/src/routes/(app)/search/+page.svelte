<h1>Song Fetcher</h1>

<script lang="ts">
  import SongDisplay from "./SongDisplay.svelte";
  import type { Song } from "$lib/types";
  import { apiFetch } from "$lib/api";
    import SongCard from "$lib/components/SongCard.svelte";

  let inputID = $state<number>(1)
  let songs = $state<Song[] | null>(null);

  async function fetchSong() {
    try {
      const data: Song = await apiFetch<Song>("/songs/" + inputID);
      console.log(data);
      songs = [data];
    } catch (error) {
      console.error(error);
    }
  }

  async function fetchAllSongs() {
    try {
      const data: Song[] = await apiFetch<Song[]>("/songs");
      songs = data;
    } catch (error) {
      console.error(error);
    }
  }

  async function deleteSong(id: number) {
    if (!songs) return;
    try {
      const data = await apiFetch("/songs/" + id, {
        method: 'DELETE',
      })
      console.log(data);
      songs = songs.filter(song => song.id != id);
    } catch (error) {
      console.error(error);
    }
  }
</script>

<div>
  <button onclick={fetchAllSongs}>Fetch All Songs</button>
</div>

<input type="number" bind:value={inputID} />

<button onclick={fetchSong}>Fetch</button>

{#if songs}
<div style="margin-top: 10px;">
  {#if songs.length == 0}
    <div>
      <span>No songs in database.</span>
    </div>
  {:else}
    <div class="song-container">
      {#each songs as song}
        <SongCard {song} />
        <button onclick={() => deleteSong(song.id)} style="margin-bottom: 10px;">Delete</button>
      {/each}
    </div>
  {/if}
</div>
{/if}



<style>
  .song-container {
    display: flex;
    align-items: flex-start;
    flex-direction: column;
  }
</style>