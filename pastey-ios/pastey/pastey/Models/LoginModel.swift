//
//  LoginModel.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-04.
//

import Foundation

struct LoginRequest: Encodable {
    let email: String
    let password: String
}

struct LoginResponse: Decodable {
    let user: LoginResponseUser
    let access_token: String
}

struct LoginResponseUser: Decodable {
    let id: Int64
    let email: String
    let ispremium: Bool
    let isemailverified: Bool
    let created_at: String
}
