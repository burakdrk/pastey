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
}
