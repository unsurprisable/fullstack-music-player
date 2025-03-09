export interface Song {
  id: number,
  title: string,
  artist: string,
  album: string,
  uploadedAt: string,
  fileURL: string,
}

export interface Playlist {
  id: number,
  name: string,
  createdAt: string,
  songIDs: number[],
}