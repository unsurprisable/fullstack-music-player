<h1>Audio Player</h1>

<script lang="ts">
    import { apiFetch } from "$lib/api";
  import type { Song } from "$lib/types";
  import SongElement from "./SongElement.svelte";
  import { onDestroy } from "svelte";

  let audio: HTMLAudioElement = new Audio();
  let playing: boolean = false;
  let progress: number = 0;
  let duration: number = 0;
  let volume: number = 1.0;
  let enableLooping: boolean = false;
  
  let songs: Song[] | null = null;
  let song: Song | null = null;
  let songIndex : number = 0;

  fetchAllSongs();

  async function fetchAllSongs() {
    try {
      const data: Song[] = await apiFetch<Song[]>("/songs");
      songs = data;
      song = songs[0];
      loadSong(song);
    } catch (error) {
      console.error(error);
      alert("Unable to fetch songs");
    }
  }

  function loadSong(song: Song) {
  if (!song || !song.fileURL) return;
    audio.src = song.fileURL;
    audio.load();

    audio.onloadedmetadata = () => {
      duration = audio.duration;
      progress = 0;

      if (playing) {
        audio.play();
      }
    };

    audio.ontimeupdate = () => {
      progress = (audio.currentTime / audio.duration) * 100;
    }

    audio.onended = () => {
      if (enableLooping) {
        audio.currentTime = 0;
        audio.play();
      } else {
        nextSong();
      }
    }
  }

  function togglePlaying() {
    if (!song) return;
    if (playing) {
      audio.pause();
    } else {
      audio.play();
    }
    playing = !playing;
  }

  function toggleLooping() {
    enableLooping = !enableLooping;
  }

  function nextSong() {
    if (!songs || songs.length === 0) return;
    songIndex = (songIndex + 1) % songs.length;
    song = songs[songIndex];
    loadSong(song);
  }

  function previousSong() {
    if (!songs || songs.length === 0) return;
    songIndex = songIndex <= 0 ? songs.length - 1 : songIndex - 1;
    song = songs[songIndex];
    loadSong(song);
  }

  function seek(event: Event) {
    const target = event.target as HTMLInputElement;
    let newTime = (Number(target.value) / 100) * audio.duration;
    audio.currentTime = newTime;
  }

  function changeVolume(event: Event) {
    const target = event.target as HTMLInputElement;
    let volume = Number(target.value);
    audio.volume = volume;
  }

  function formatTime(seconds: number) {
    const h = Math.floor(seconds / 3600);
    const m = Math.floor((seconds % 3600) / 60);
    const s = Math.floor(seconds % 60);

    if (h > 0) {
      return `${h}:${m.toString().padStart(2, "0")}:${s.toString().padStart(2, "0")}`;
    } else {
      return `${m}:${s.toString().padStart(2, "0")}`;
    }
  }

  onDestroy(() => {
    audio.pause();
    audio.currentTime = 0;
  })
</script>

{#if songs}
  <div>
    <h3>Choose a song:</h3>
    {#each songs as s, i}
     <div style="padding: 2px;">
      <SongElement song={s} onclick={() => {
        song = s;
        songIndex = i;
        playing = true;
        loadSong(song)
      }} />
     </div>
    {/each}
  </div>
{:else}
  <div>
    <p>Please select a song.</p>
  </div>
{/if}

{#if song}
  <h3>Now playing:</h3>
  <p>{song.artist} {song.artist != "" ? "-" : ""} {song.title} {song.album != "" ? "(" + song.album + ")" : ""}</p>
{/if}

<div style="padding-top: 10px;">
  <button onclick={togglePlaying}>{playing ? "Pause" : "Play"}</button>
  <input type="range" min="0" max="100" step="0.1" bind:value={progress} oninput={seek} style="width:200px"/>
  <span>{formatTime(progress / 100 * duration)} / {formatTime(duration)}</span>
  <input type="range" min="0" max="1" step="0.01" bind:value={volume} oninput={changeVolume} style="width: 80px;" />
  <div>
    <button onclick={previousSong}>Previous</button>
    <button onclick={nextSong}>Next</button>
    <button onclick={toggleLooping}>Loop ({enableLooping ? "on" : "off"})</button>
  </div>
</div>