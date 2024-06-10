//
//  EncryptionService.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-06.
//

import Foundation
import SwiftyRSA
import SimpleKeychain

class EncryptionService {
    static let shared = EncryptionService()
    private let keychain = SimpleKeychain()
    
    private init() {}
    
    func generateKeyPair(keySize: Int = 2048) throws -> (publicKey: String, privateKey: String) {
        let keyPair = try SwiftyRSA.generateRSAKeyPair(sizeInBits: keySize)
        
        let publicKey = try keyPair.publicKey.pemString()
        let privateKey = try keyPair.privateKey.pemString()
        
        try keychain.set(privateKey, forKey: "privateKey")
        try keychain.set(publicKey, forKey: "publicKey")
        
        return (publicKey, privateKey)
    }
    
    func encryptMessage(message: String, publicKey: String) throws -> String {
        let key = try PublicKey(pemEncoded: publicKey)
        let clear = try ClearMessage(string: message, using: .utf8)
        let encrypted = try clear.encrypted(with: key, padding: .PKCS1)
        
        return encrypted.base64String
    }
    
    func decryptMessage(encryptedMessage: String) throws -> String {
        let key = try PrivateKey(pemEncoded: try keychain.string(forKey: "privateKey"))
        let encrypted = try EncryptedMessage(base64Encoded: encryptedMessage)
        let decrypted = try encrypted.decrypted(with: key, padding: .PKCS1)
        
        return try decrypted.string(encoding: .utf8)
    }
}
