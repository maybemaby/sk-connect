import {  sampleClient } from "$lib/api/api.js";
import type { PageLoad } from "./$types.js";

export const load: PageLoad = async (event) => {
    const name = event.url.searchParams.get("name") ?? "world";


    const res = await sampleClient.sampleMethod({
        name,
    })


    return {
        msg: res.message
    }
}