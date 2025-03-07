<h1>Upload Song</h1>

<script lang="ts">
    import { apiFetch } from "$lib/api";

  let files = $state<FileList | null>(null);

  function handleFileChange(event: Event) {
    const input = event.target as HTMLInputElement;
    if (input && input.files) {
      files = input.files;
    }
  }

  async function uploadFiles() {
    if (!files || files.length === 0) {
      alert('Please select at least 1 file to upload!');
      return;
    }

    const formData = new FormData();
    for (let i = 0; i < files.length; i++) {
      formData.append('files', files[i])
    }

    try {
      const data = await apiFetch("/upload", {
        method: 'POST',
        body: formData
      });
      console.log(data);
    } catch (error) {
      console.error(error);
    }
  }
</script>

<input type="file" accept=".mp3" multiple onchange={handleFileChange}>

<div style="margin-top: 20px;">
  <button onclick={uploadFiles} class="fancy-button">
    <span>Upload</span>
  </button>
</div>

<style>
  .fancy-button {
    background: linear-gradient(135deg, #6e8efb, #a777e3);
    color: white;
    border: none;
    border-radius: 8px;
    padding: 12px 24px;
    font-size: 16px;
    font-weight: bold;
    cursor: pointer;
    transition: all 0.3s ease;
    box-shadow: 0 4px 6px rgba(50, 50, 93, 0.11), 0 1px 3px rgba(0, 0, 0, 0.08);
    outline: none;
    position: relative;
    overflow: hidden;
  }

  .fancy-button:hover {
    transform: translateY(-2px);
    box-shadow: 0 7px 14px rgba(50, 50, 93, 0.1), 0 3px 6px rgba(0, 0, 0, 0.08);
  }

  .fancy-button:active {
    transform: translateY(1px);
  }
</style>