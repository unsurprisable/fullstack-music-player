<script>
  let { href, children } = $props();
</script>

<h3>
  <a {href} style="font-size: xx-large; font-weight: bolder;">
    {@render children()}
  </a>
</h3>