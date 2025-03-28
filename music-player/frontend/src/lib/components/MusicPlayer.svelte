<script lang="ts">
  import { changeVolume, currentSong, duration, enableLooping, enableShuffle, isPlaying, nextSong, previousSong, progress, seek, toggleLooping, togglePlaying, toggleShuffle, volume } from "$lib/stores/player";
  import { Play, Pause, Volume2, Shuffle, SkipBack, Repeat, SkipForward } from "lucide-svelte";
  
  function formatTime(seconds: number) {
    if (isNaN(seconds)) return "0:00";
    const h = Math.floor(seconds / 3600);
    const m = Math.floor((seconds % 3600) / 60);
    const s = Math.floor(seconds % 60);

    if (h > 0) {
      return `${h}:${m.toString().padStart(2, "0")}:${s.toString().padStart(2, "0")}`;
    } else {
      return `${m}:${s.toString().padStart(2, "0")}`;
    }
  }
</script>

<div class="music-player">
  <!-- Song Info -->
  <div class="left-section">
    <div class="song-info">
      {#if $currentSong}
        <span class="song-title">{$currentSong.title}</span>
        <span class="text">{$currentSong.artist}</span>
      {:else}
        <span class="song-title">Song Title</span>
        <span class="text">Artist</span>
      {/if}
    </div>
  </div>

  <!-- Playback Controls -->
  <div class="playback-container">
    <div class="playback-controls">
      <div class="control-button-container">
        <button class="control-button" onclick={toggleShuffle}>
          <Shuffle size="18" color={$enableShuffle ? "lime" : "white"}/>
          {#if $enableShuffle}
            <span class="indicator-dot"></span>
          {/if}
        </button>
      </div>
      <button class="control-button" onclick={previousSong}><SkipBack size="22" /></button>
      <button class="play-button" onclick={togglePlaying}>{#if $isPlaying}<Pause size="22" fill="true" />{:else}<Play size="22" fill="true"/>{/if}</button>
      <button class="control-button" onclick={nextSong}><SkipForward size="22" /></button>
      <div class="control-button-container">
        <button class="control-button" onclick={toggleLooping}>
          <Repeat size="18" color= {$enableLooping ? "lime" : "white"}/>
          {#if $enableLooping}
            <spon class="indicator-dot"></spon>
          {/if}
        </button>
      </div>
    </div>
    <div class="progress-bar">
      <span class="text">{formatTime($progress * $duration / 100)}</span>
      <input type="range" min="0" max="100" step="0.1" bind:value={$progress} oninput={(e) => seek(Number((e.target as HTMLInputElement).value))} style="--percent: {$progress}%;" />
      <span class="text">{formatTime($duration)}</span>
    </div>
  </div>

  <!-- Extra Controls -->
  <div class="right-section">
    <div class="extra-controls">
      <span><Volume2 size="20" /></span>
      <input type="range" min="0" max="1" step="0.025" bind:value={$volume} oninput={(e) => changeVolume(Number((e.target as HTMLInputElement).value))} style="--percent: {$volume*100}%;" />
    </div>
  </div>
</div>






<style>
  .music-player {
    font-family: 'Gill Sans', 'Gill Sans MT', Calibri, 'Trebuchet MS', sans-serif;
    position: fixed;
    bottom: 0;
    left: 0;
    width: 100%;
    background: #050505;
    color: white;
    display: flex;
    align-items: center;
    justify-content: space-between;
    box-sizing: border-box;
    padding: 20px 20px 16px 20px;
    box-shadow: -50px -5px 12px rgba(0, 0, 0, .25);
  }
  .music-player button, input {
    cursor: pointer;
  }
  .music-player input[type="range"] {
    appearance: none;
    height: 5px;
    border-radius: 2.5px;
    background: linear-gradient(to right,
      white 0%,
      white var(--percent, 0%),
      #3f3f3f var(--percent, 0%),
      #3f3f3f 100%);
  }
  .music-player input[type="range"]::-webkit-slider-thumb {
    appearance: none;
    width: 12px;
    height: 12px;
    background: white;
    border-radius: 50%;
  }
  .left-section, .right-section {
    display: flex;
    align-items: center;
    flex: 1;
  }
  .left-section {
    justify-content: flex-start;
    flex: 1;
  }
  .right-section {
    justify-content: flex-end;
    flex: 1;
  }
  .song-info {
    display: flex;
    align-items: flex-start;
    flex-direction: column;
    gap: 8px;
  }
  .playback-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    flex: 2;
  }
  .playback-controls {
    display: flex;
    align-items: center;
    flex-direction: row;
    gap: 20px;
    padding-left: 19px;
    padding-bottom: 10px;
  }
  .progress-bar {
    display: flex;
    align-items: center;
    flex-direction: row;
    justify-content: center;
    width: 100%;
  }
  .progress-bar input[type="range"] {
    width: 100%;
    margin: 0 10px;
  }
  .play-button {
    width: 35px;
    height: 35px;
    border-radius: 50%;
    border: none;
    margin-top: -2px;
  }
  .control-button-container {
    position: relative;
    display: flex;
    justify-content: center;
    align-items: center;
  }
  .control-button {
    background: none;
    color: white;
    border: none;
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0;
  }
  .indicator-dot {
    position: absolute;
    width: 4px;
    height: 4px;
    background-color: lime;
    border-radius: 50%;
    bottom: -8px;
    left: 50%;
    transform: translateX(-50%);
  }
  .extra-controls {
    display: flex;
    align-items: center;
    width: 50%
  }
  .extra-controls input[type="range"] {
    width: 100%;
    margin-left: 8px;
  }
  .song-title {
    font-size: 18px;
    font-weight: 600;
  }
  .text {
    font-size: 14px;
    font-weight: 500;
    color: #cfcfcf;
  }

</style>