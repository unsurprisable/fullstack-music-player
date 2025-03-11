import { browser } from "$app/environment";
import { apiFetch } from "$lib/api";
import type { Playlist, Song } from "$lib/types";
import { get, writable } from "svelte/store";


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
export const enableShuffle = writable(false);

let shuffleArray: number[] | null = null;

let audio: HTMLAudioElement;

if (browser) {
  audio = new Audio();
}

export function loadSong(song: Song, autoPlay: boolean = false) {
  if (!song || !song.fileURL) return;
  audio.src = song.fileURL;
  audio.load();
  currentSong.set(song);

  if (autoPlay) {
    isPlaying.set(true);
  }

  audio.onloadedmetadata = () => {
    duration.set(audio.duration);
    isPlaying.update(p => {
      if (p) audio.play();
      return p;
    });
  };

  audio.ontimeupdate = () => {
    progress.set((audio.currentTime / audio.duration) * 100);
  }

  audio.onended = () => {
    if (!get(currentPlaylist)) {
      if (get(enableLooping)) {
        audio.currentTime = 0;
        audio.play();
      } else {
        isPlaying.set(false);
      }
    } else {
      currentSongs.update(songs => {
        if (!songs || songs.length === 0) return null;
        const index = songs.findIndex(s => s.id === song?.id);
        if (index === songs.length-1 && !get(enableLooping)) {
          isPlaying.set(false);
        } else {
          nextSong();
        }
        return songs;
      });
    }
  };
}

function playlistsEqual(playlist1: Playlist | null, playlist2: Playlist | null): boolean {
  if (playlist1 === playlist2) return true;
  if (!playlist1 || !playlist2) return false;
  return playlist1.id === playlist2.id && playlist1.name === playlist2.name;
}

export function setPlaylist(playlist: Playlist | null, songs: Song[] | null = null) {
  if (playlist && !songs) {
    console.error("Tried to set playlist without any songs!");
    return;
  }

  if (!playlistsEqual(playlist, get(currentPlaylist))) {
    currentPlaylist.set(playlist);
    currentSongs.set(songs);
    console.log("Set new playlist to: " + JSON.stringify(playlist));
  }

  if (playlist && get(enableShuffle)) {
    console.log("shuffle time")
    generateShuffleArray();
  } else if (!playlist) {
    shuffleArray = null;
  }
}

export function togglePlaying() {
  isPlaying.update(p => {
    if (p) {
      audio.pause();
    } else {
      audio.play();
    }
    return !p;
  });
}

export function toggleLooping() {
  enableLooping.update(e => {
    return !e;
  });
}

export function toggleShuffle() {
  enableShuffle.update(e => {
    if (e) {
      shuffleArray = null;
    } else if (get(currentSongs)) {
      generateShuffleArray();
    }
    return !e;
  });
}

function generateShuffleArray(ignoreCurrentSong: boolean = false) {
  if (!currentSongs || !get(currentSongs)) return;
  const songs : Song[] | null = get(currentSongs);
  const songAmount: number = songs?.length ?? 0;

  // yeah i wont even try to lie and say i did this part myself lol (i have no idea whats happening in these lines)
  shuffleArray = Array.from({ length: songAmount }, (_, i) => i);
  for (let i = shuffleArray.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [shuffleArray[i], shuffleArray[j]] = [shuffleArray[j], shuffleArray[i]];
  }

  if (!currentSong || !get(currentSongs) || ignoreCurrentSong) { 
    console.log(shuffleArray); 
    return; 
  }

  const index = songs?.findIndex(s => s.id === get(currentSong)?.id);
  shuffleArray = shuffleArray.filter(i => i !== index);
  console.log(shuffleArray); 
}

export function nextSong() {
  if (!get(currentSongs) || get(currentSongs)?.length === 0) {
    audio.currentTime = 0;
    audio.play();
    return;
  }
  currentSongs.update(songs => {
    currentSong.update(song => {
      if (!songs) return null;

      if (get(enableShuffle) && shuffleArray) {
        const nextIndex = shuffleArray.pop() ?? 0;
        if (shuffleArray.length === 0) {
          generateShuffleArray(true);
        } else {
          console.log(shuffleArray);
        }
        loadSong(songs[nextIndex]);
        return songs[nextIndex];
      } else {
        const index = songs.findIndex(s => s.id === song?.id);
        const nextIndex = (index + 1) % songs.length;
        loadSong(songs[nextIndex]);
        return songs[nextIndex];
      }
    });
    return songs;
  });
}

export function previousSong() {
  // for now this doesnt work when shuffle is on because that sounds like a whole lotta work that i dont feel like doing
  if (!get(currentSongs) || get(currentSongs)?.length === 0 || get(enableShuffle)) {
    audio.currentTime = 0;
    audio.play();
    return;
  }
  currentSongs.update(songs => {
    currentSong.update(song => {
      if (!songs) return null;
      const index = songs.findIndex(s => s.id === song?.id);
      const prevIndex = index <= 0 ? songs.length - 1 : index - 1;
      loadSong(songs[prevIndex]);
      return songs[prevIndex];
    });
    return songs;
  });
}

export function seek(newProgress: number) {
  audio.currentTime = (newProgress / 100) * audio.duration;
  if (!get(isPlaying)) {
    togglePlaying();
  }
}

export function changeVolume(newVolume: number) {
  audio.volume = newVolume;
}