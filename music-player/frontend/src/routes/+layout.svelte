<!-- MAIN LAYOUT FOR THE MUSIC PLAYER -->
 
<script>
  import MusicPlayer from "$lib/components/MusicPlayer.svelte";
</script>

<div class="main-content">
  <slot />
</div>

<MusicPlayer />

<style>
  /* so the music player doesn't block stuff at the bottom and u can scroll further up */
  .main-content {
    margin-bottom: 140px;
  }
</style>