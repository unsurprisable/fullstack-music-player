<h1>Audio Player</h1>

<script lang="ts">
    import { apiFetch } from "$lib/api";
  import type { Playlist, Song } from "$lib/types";
    import PlaylistElement from "./PlaylistElement.svelte";
  import SongElement from "./SongElement.svelte";
  import { onDestroy, onMount } from "svelte";

  let audio: HTMLAudioElement = new Audio();
  let playing: boolean = false;
  let progress: number = 0;
  let duration: number = 0;
  let volume: number = 1.0;
  let enableLooping: boolean = false;
  
  let playlists: Playlist[] | null = null;
  let playlist: Playlist | null = null
  let songs: Song[] | null = null;
  let song: Song | null = null;
  let songIndex : number = 0;

  async function fetchAllPlaylists() {
    try {
      const data: Playlist[] = await apiFetch<Playlist[]>("/playlists");
      playlists = data;
    } catch (error) {
      console.error(error);
      alert("Unable to fetch songs");
    }
  }

  async function fetchSongsFromPlaylist() {
    if (!playlist) return;
    try {
      const data: Song[] = await apiFetch<Song[]>("/playlists/"+playlist.id+"/songs");
      console.log(data);
      songs = data;
      song = songs[0]
      songIndex = 0;
      playing = true;
      loadSong(song)
    } catch (error) {
      console.error(error);
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

  onMount(() => {
    fetchAllPlaylists();
  })
</script>

{#if playlists}
  <div>
    <h3>Choose a playlist:</h3>
    {#each playlists.toReversed() as p}
     <div style="padding: 2px;">
      <PlaylistElement playlist={p} onclick={() => {
        playlist = p;
        fetchSongsFromPlaylist();
      }} />
     </div>
    {/each}
  </div>
{:else}
  <div>
    <p>Please select a playlist.</p>
  </div>
{/if}

{#if songs && playlist}
  <div>
    <h3>Tracklist: {playlist.name}</h3>
    {#each songs as s, i}
      <div>
        {#if songIndex == i}
          <span>{"→"} </span>
        {/if}
        <SongElement song={s} onclick={() => {
          songIndex = i
          song = s
          loadSong(s)
        }} />
      </div>
    {/each}
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