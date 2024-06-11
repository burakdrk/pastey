//
//  DeveloperPreview.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-11.
//

import Foundation

class DeveloperPreview {
    static let instance = DeveloperPreview()
    private init() {}
    
    let entry = Entry(
        id: 1,
        entryID: UUID().uuidString,
        userID: 3,
        fromDeviceID: 3,
        toDeviceID: 3,
        encryptedData: "encryptedDataencryptedDataencryptedDataencryptedDataencryptedDataencryptedData",
        createdAt: "createdAt",
        fromDeviceName: "Windows"
    )
}
