<h1>Admin Page</h1>

<script>
    import { apiFetch } from "$lib/api";

  async function clearDatabase() {
    try {
      const data = await apiFetch("/songs/clear");
      console.log(data)
    } catch (error) {
      console.error(error)
    }
  }
</script>

<button onclick={clearDatabase}>Clear Database</button>