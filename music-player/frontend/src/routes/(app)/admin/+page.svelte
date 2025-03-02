<h1>Admin Page</h1>

<script>
  async function clearDatabase() {
    try {
      const response = await fetch("http://localhost:8080/songs/clear")
      if (!response.ok) throw new Error();

      const data = await response.json();
      console.log(data);
    } catch (error) {
      console.error("Error cleaing database: ", error);
    }
  }
</script>

<button onclick={clearDatabase}>Clear Database</button>