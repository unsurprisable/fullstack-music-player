<script lang="ts">
  import { loadSong, setPlaylist } from "$lib/stores/player";
  import type { Playlist, Song } from "$lib/types";
  import { Play } from "lucide-svelte";

  export let song: Song | null = null;
  export let onPlay: (() => void) | null = null;
  export let inPlaylist: boolean = false;

  function onclick() {
    if (!song) return;
    loadSong(song, true);
    if (onPlay) onPlay();
    if (!inPlaylist) setPlaylist(null);
  }
</script>

<div class="container">
  <div class="play-button">
    <button {onclick}><Play size="22" fill="true" opacity={song ? 100 : 0}/></button>
  </div>
  <div class="song-info">
    <span class="song-title">{song ? song.title : "Title"}</span>
    <span class="text">{song ? song.artist : ""}</span>
  </div>
  <div class="album">
    <span class="text">{song ? song.album : "Album"}</span>
  </div>
  <div class="upload-date">
    <span class="text">{song ? song.uploadedAt.split("T")[0] : "Date added"}</span>
  </div>
</div>



<style>
  .container {
    font-family: 'Gill Sans', 'Gill Sans MT', Calibri, 'Trebuchet MS', sans-serif;
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 800px;
    height: 60px;
    background: #050505;
    color: white;
  }
  .text {
    color: #afafaf;
    font-weight: 500;
  }
  .song-title {
    color: white;
    font-weight: 600;
    font-size: 18px;
  }
  .play-button {
    display: flex;
    padding: 0 20px 0 10px;
    flex: 0;
  }
  .play-button button {
    color: white;
    background: none;
    border: none;
  }
  .song-info {
    display: flex;
    align-items: flex-start;
    flex-direction: column;
    gap: 2px;
    flex: 2;
  }
  .album {
    display: flex;
    align-items: flex-start;
    flex: 2;
  }
  .upload-date {
    display: flex;
    align-items: flex-start;
    flex: 1
  }
</style>