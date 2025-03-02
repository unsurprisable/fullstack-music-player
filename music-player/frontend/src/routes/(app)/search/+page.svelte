<h1>Song Fetcher</h1>

<script lang="ts">
  import SongDisplay from "./SongDisplay.svelte";
  import type { Song } from "$lib/types";

  let id = $state<number>(1)
  let songs = $state<Song[] | null>(null);

  async function fetchSong() {
    try {
      const response = await fetch("http://localhost:8080/songs/" + id);
      if (!response.ok) throw new Error("Failed to fetch song (response not OK)");
      
      const data: Song = await response.json();
      console.log(data);
      songs = [data];
    } catch (error) {
      console.error("Error fetching song: ", error);
      alert("Unable to fetch songs");
    }
  }
  async function fetchAllSongs() {
    try {
      const response = await fetch("http://localhost:8080/songs");
      if (!response.ok) throw new Error("Failed to fetch songs (response not OK)");
      
      const data: Song[] = await response.json();
      console.log(data);
      songs = data;
    } catch (error) {
      console.error("Error fetching songs: ", error);
      alert("Unable to fetch songs");
    }
  }
</script>

<div>
  <button onclick={fetchAllSongs}>Fetch All Songs</button>
</div>

<input type="number" bind:value={id} />

<button onclick={fetchSong}>Fetch</button>

{#if songs}
  {#each songs as song}
    <SongDisplay {song} />
  {/each}
{/if}