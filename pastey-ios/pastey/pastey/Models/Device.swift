//
//  Device.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-05.
//

import Foundation

// MARK: - Device
struct Device: Codable, Identifiable {
    let id, userID: Int
    let deviceName, publicKey, createdAt: String

    enum CodingKeys: String, CodingKey {
        case id
        case userID = "user_id"
        case deviceName = "device_name"
        case publicKey = "public_key"
        case createdAt = "created_at"
    }
}

typealias ListDevicesResponse = [Device]

// MARK: - CreateDeviceRequest
struct CreateDeviceRequest: Codable {
    let deviceName, publicKey: String

    enum CodingKeys: String, CodingKey {
        case deviceName = "device_name"
        case publicKey = "public_key"
    }
}

// MARK: - CreateDeviceResponse
struct CreateDeviceResponse: Codable {
    let device: Device
    let sessionID, refreshToken, refreshTokenExpiresAt: String

    enum CodingKeys: String, CodingKey {
        case device
        case sessionID = "session_id"
        case refreshToken = "refresh_token"
        case refreshTokenExpiresAt = "refresh_token_expires_at"
    }
}
