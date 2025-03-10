import type { Playlist, Song } from "$lib/types";
import { writable } from "svelte/store";


// song info
export const currentPlaylist = writable<Playlist | null>(null);
export const currentSongs = writable<Song[] | null>(null);
export const currentSong = writable<Song | null>(null);

// player states
export const isPlaying = writable(false);
export const progress = writable(0);
export const duration = writable(0);
export const volume = writable(1.0);
export const enableLooping = writable(false);

let audio: HTMLAudioElement = new Audio();

export function loadSong(song: Song) {
  if (!song || !song.fileURL) return;
  audio.src = song.fileURL;
  audio.load();

  audio.onloadedmetadata = () => {
    duration.set(audio.duration);
    progress.set(0);
    isPlaying.update(p => {
      if (p) audio.play();
      return p;
    });
  };

  audio.ontimeupdate = () => {
    progress.set((audio.currentTime / audio.duration) * 100);
  }

  audio.onended = () => {
    if (enableLooping) {
      audio.currentTime = 0;
      audio.play();
    } else {
      // next song
    }
  };

  currentSong.set(song);
}

export function togglePlaying() {
  // todo
}