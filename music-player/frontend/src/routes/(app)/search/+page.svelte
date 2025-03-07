<h1>Song Fetcher</h1>

<script lang="ts">
  import SongDisplay from "./SongDisplay.svelte";
  import type { Song } from "$lib/types";
  import { apiFetch } from "$lib/api";

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
      console.log(data);
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
  {#if songs.length == 0}
    <div>
      <span>No songs in database.</span>
    </div>
  {:else}
    {#each songs as song}
      <SongDisplay {song} onclick={() => {deleteSong(song.id)}} />
    {/each}
  {/if}
{/if}