import {createClient} from "@connectrpc/connect"
import { createConnectTransport } from "@connectrpc/connect-web"
import {SampleService} from "$lib/gen/proto/api/v1/sample_connect.js"

const transport = createConnectTransport({
    baseUrl: "http://localhost:8000",
})

export const sampleClient = createClient(SampleService, transport)
