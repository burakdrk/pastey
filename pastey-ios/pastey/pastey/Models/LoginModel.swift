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
    let user: User
    let access_token: String
}
