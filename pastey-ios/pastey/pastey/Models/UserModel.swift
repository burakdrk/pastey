//
//  UserModel.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-04.
//

import Foundation

struct User: Decodable {
    let id: Int64
    let email: String
    let ispremium: Bool
    let isemailverified: Bool
    let created_at: String
}
