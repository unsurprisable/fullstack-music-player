<script>

    import MusicPlayer from "$lib/components/MusicPlayer.svelte";

</script>
<slot />

<div class="player-container">
  <MusicPlayer />
</div>

<style>
  .player-container {
    position: fixed;
    bottom: 0;
    left: 0;
    width: 100%;
    background: #222;
    color: white;
    padding: 30px;
    box-shadow: 0 -2px 5px rgba(0, 0, 0, 0.2);
  }
</style>