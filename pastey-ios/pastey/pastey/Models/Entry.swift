//
//  Entry.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-11.
//

import Foundation

// MARK: - Entry
struct Entry: Codable, Identifiable {
    let id: Int
    let entryID: String
    let userID, fromDeviceID, toDeviceID: Int
    let encryptedData, createdAt, fromDeviceName: String

    enum CodingKeys: String, CodingKey {
        case id
        case entryID = "entry_id"
        case userID = "user_id"
        case fromDeviceID = "from_device_id"
        case toDeviceID = "to_device_id"
        case encryptedData = "encrypted_data"
        case createdAt = "created_at"
        case fromDeviceName = "from_device_name"
    }
}

typealias EntryResponse = [Entry]


// MARK: - CopyRequest
struct CopyRequest: Codable {
    let fromDeviceID: Int
    let copies: [Copy]

    enum CodingKeys: String, CodingKey {
        case fromDeviceID = "from_device_id"
        case copies
    }
}

// MARK: - Copy
struct Copy: Codable {
    let toDeviceID: Int
    let encryptedData: String

    enum CodingKeys: String, CodingKey {
        case toDeviceID = "to_device_id"
        case encryptedData = "encrypted_data"
    }
}
