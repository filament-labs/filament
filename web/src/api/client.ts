import { PingService } from "@/pb/ping_pb"
import { WalletService } from "@/pb/wallet_pb"
import { createClient } from "@connectrpc/connect"
import { createConnectTransport } from "@connectrpc/connect-web"


const transport = createConnectTransport({
  baseUrl: "http://localhost:8080"
})

export const pingClient = createClient(PingService, transport)
export const walletClient = createClient(WalletService, transport)