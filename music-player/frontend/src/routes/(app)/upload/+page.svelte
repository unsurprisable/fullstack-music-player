<h1>Upload Song</h1>

<script lang="ts">
  let file = $state<File | null>(null);

  function handleFileChange(event: Event) {
    const input = event.target as HTMLInputElement;
    if (input && input.files) {
      file = input.files[0];
    }
  }

  async function uploadFile() {
    if (!file) {
      alert('Please select a file to upload!');
      return;
    }

    const formData = new FormData();
    formData.append('file', file);

    try {
      const response = await fetch('http://localhost:8080/upload', {
        method: 'POST',
        body: formData,
      });

      if (!response.ok) {
        throw new Error('Failed to upload file');
      }

      const data = await response.json();
      console.log(data);
    } catch (error) {
      console.error("Error uploading file: ", error);
    }
  }
</script>

<input type="file" onchange={handleFileChange}>

<div style="margin-top: 20px;">
  <button onclick={uploadFile} class="fancy-button">
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