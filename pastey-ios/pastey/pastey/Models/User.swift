//
//  User.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-04.
//

import Foundation

// MARK: - User
struct User: Codable, Identifiable {
    let id: Int
    let email: String
    let ispremium, isemailverified: Bool
    let createdAt: String

    enum CodingKeys: String, CodingKey {
        case id, email, ispremium, isemailverified
        case createdAt = "created_at"
    }
}

// MARK: - CreateUserRequest
struct CreateUserRequest: Codable {
    let email, password: String
}

// MARK: - LoginRequest
struct LoginRequest: Codable {
    let email, password: String
    let deviceID: Int?

    enum CodingKeys: String, CodingKey {
        case email, password
        case deviceID = "device_id"
    }
}

// MARK: - LoginResponse
struct LoginResponse: Codable {
    let accessToken, accessTokenExpiresAt: String
    let sessionID, refreshToken, refreshTokenExpiresAt: String?
    let user: User

    enum CodingKeys: String, CodingKey {
        case sessionID = "session_id"
        case accessToken = "access_token"
        case accessTokenExpiresAt = "access_token_expires_at"
        case refreshToken = "refresh_token"
        case refreshTokenExpiresAt = "refresh_token_expires_at"
        case user
    }
}

// MARK: - RefreshTokenRequest
struct RefreshTokenRequest: Codable {
    let refreshToken: String

    enum CodingKeys: String, CodingKey {
        case refreshToken = "refresh_token"
    }
}

// MARK: - RefreshTokenResponse
struct RefreshTokenResponse: Codable {
    let accessToken, accessTokenExpiresAt: String

    enum CodingKeys: String, CodingKey {
        case accessToken = "access_token"
        case accessTokenExpiresAt = "access_token_expires_at"
    }
}
