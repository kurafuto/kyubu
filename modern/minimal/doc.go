// Package "minimal" (actually "modern" for compatibility reasons) implements a
// subset of the Minecraft protocol. It basically contains enough for Kurafuto
// to work, as well as a generic packet that acts as a catch-all and just
// throws around bytes.
//
// Packets we care about:
// >> Handshake
//  * 0x00 Handshake
// >> Login (Clientbound)
//  * 0x00 LoginDisconnect
//  * 0x01 EncryptionRequest
//  * 0x02 LoginSuccess
//  * 0x03 SetInitialCompression
// >> Login (Serverbound)
//  * 0x00 LoginStart
//  * 0x01 EncryptionResponse
// >> Play (Clientbound)
//  * 0x01 JoinGame
//  * 0x02 ServerMessage
//  * 0x07 Respawn
//  * 0x38 PlayerListItem
//  * 0x3A ServerTabComplete
//  * 0x3B ScoreboardObjective
//  * 0x3C UpdateScore
//  * 0x3D ShowScoreboard
//  * 0x3E Teams
//  * 0x3G ServerPluginMessage
//  * 0x40 Disconnect
//  * 0x46 SetCompression
// >> Play (Serverbound)
//  * 0x01 ClientMessage
//  * 0x14 ClientTabComplete
//  * 0x16 ClientStatus
//  * 0x17 ClientPluginMessage
// >> Status (Clientbound)
//  * 0x00 StatusResponse
//  * 0x01 StatusPong
// >> Status (Serverbound)
//  * 0x00 StatusRequest
//  * 0x01 StatusPing
package modern
