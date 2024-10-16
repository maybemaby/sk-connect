<script lang="ts">
    import "../app.css";

    let {
        data
    } = $props();

</script>


<h1>SK-Connect</h1>

<div>
    <p>Server Response:</p>
    <code>
        {JSON.stringify(data, null, 2)}
    </code>
</div>

<style>

</style>