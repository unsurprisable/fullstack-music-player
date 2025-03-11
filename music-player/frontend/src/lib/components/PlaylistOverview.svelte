<script lang="ts">
    import { apiFetch } from "$lib/api";
    import type { Playlist, Song } from "$lib/types";
    import { onMount } from "svelte";
    import SongCard from "./SongCard.svelte";
    import { setPlaylist } from "$lib/stores/player";

  export let playlist: Playlist | null = null;
  let songs: Song[] | null = null;
  
  async function fetchSongs() {
    if (!playlist) return;
    try {
      let data: Song[] = await apiFetch<Song[]>("/playlists/"+ playlist.id +"/songs");
      songs = data;
    } catch (error) {
      console.error(error);
    }
  }

  onMount(() => {
    fetchSongs();
  })
</script>



<div class="playlist">
  <div class="playlist-info">
    <span class="playlist-name">{playlist?.name}</span>
    <span class="song-count">{songs?.length} songs</span>
  </div>
  <div class="song-container">
    {#if songs}
      <SongCard />
      <div style="width: 100%; height: 1px; background: white; opacity: 30%;"></div>
      {#each songs as song}
        <SongCard {song} inPlaylist={true} onPlay={() => setPlaylist(playlist, songs)} />
      {/each}
      <div style="width: 100%; height: 1px; background: white; opacity: 30%;"></div>
    {:else}
      <span>No songs {":("}</span>
    {/if}
  </div>
</div>




<style>
  .playlist {
    display: flex;
    flex-direction: column;
    width: 1000px;
    height: 100%;
    background: #050505;
    color: white;
    padding: 0 50px;
  }
  .playlist-info {
    display: flex;
    flex-direction: column;
    padding: 25px 20px;
  }
  .playlist-name {
    font-size: 80px;
    font-weight: 900;
  }
  .song-count {
    font-size: 24px;
  }
  .song-container {
    display: flex;
    flex-direction: column;
    padding-bottom: 50px;
  }
</style>