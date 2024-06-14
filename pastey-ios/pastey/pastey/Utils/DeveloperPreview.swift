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
        createdAt: "2024-06-11T06:24:41.805469Z",
        fromDeviceName: "Windows"
    )
    
    let user = User(id: 1, email: "test@tset.tset", ispremium: true, isemailverified: true, createdAt: "2024-06-11T06:24:41.805469Z")
    
    func createEntryArray() -> [Entry] {
        var entries = [Entry]()
        
        for i in 1...10 {
            let temp = Entry(
                id: i,
                entryID: UUID().uuidString,
                userID: 3,
                fromDeviceID: 3,
                toDeviceID: 3,
                encryptedData: "\(i)encryptedData",
                createdAt: "2024-06-11T06:24:41.805469Z",
                fromDeviceName: "Windows"
            )
            
            entries.append(temp)
        }
        
        return entries
    }
    
    func createDeviceArray() -> [Device] {
        var devices = [Device]()
        
        for i in 1...5 {
            let temp = Device(
                id: i,
                userID: 3,
                deviceName: "Windows",
                publicKey: "publicKey",
                createdAt: "2024-06-11T06:24:41.805469Z"
            )
            
            devices.append(temp)
        }
        
        return devices
    }
}
