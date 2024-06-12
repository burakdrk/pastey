//
//  HomeViewModel.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-11.
//

import Foundation
import SwiftUI

class HomeViewModel: ObservableObject {
    @Published var entries: [Entry]
    @Published var editMode = EditMode.inactive
    @Published var errorMessage: String?
    @Published var isFetching = false
    
    @AppStorage("deviceID") var deviceID: Int = -1
    
    private let entryService: EntryService
    private let deviceService: DeviceService

    init() {
        self.entryService = EntryService()
        self.deviceService = DeviceService()
        self.entries = DeveloperPreview.instance.createEntryArray()
        
        if self.deviceID == -1 { return }
        
        Task { await fetchEntries() }
    }
    
    func deleteEntry(at offsets: IndexSet) {
        isFetching = true
        
        Task { @MainActor in
            let entry = entries[offsets.first!]
            let err = await entryService.deleteEntry(entryID: entry.entryID)
            
            errorMessage = err
            if err != nil {
                isFetching = false
                return
            }
            
            self.entries.remove(atOffsets: offsets)
            isFetching = false
        }
    }
}

// MARK: - Paste
extension HomeViewModel {
    func paste() {
        let clipboardString = ClipboardService.shared.pasteFromClipboard()
        
        guard let clipboardString else {
            return
        }
        
        isFetching = true
                
        Task { @MainActor in
            var (res, err) = await deviceService.fetchDevices()
            
            errorMessage = err
            guard let res else {
                isFetching = false
                return
            }
            
            var copies = [Copy]()
                        
            for device in res {
                let encrypted = try? EncryptionService.shared.encryptMessage(message: clipboardString, publicKey: device.publicKey)
                
                guard let encrypted else {
                    errorMessage = "Encryption failed"
                    isFetching = false
                    return
                }
                
                copies.append(Copy(toDeviceID: device.id, encryptedData: encrypted))
            }
            
            err = await entryService.copyEntries(fromDeviceID: self.deviceID, copies: copies)
            
            if let err {
                errorMessage = err
                isFetching = false
                return
            }
            
            await fetchEntries()
        }
    }
}

// MARK: - Copy
extension HomeViewModel {
    func copy() {
        isFetching = true
        
        Task { @MainActor in
            await fetchEntries()
            if self.entries.isEmpty {
                isFetching = false
                return
            }
            
            let decrypted = try? EncryptionService.shared.decryptMessage(encryptedMessage: self.entries[0].encryptedData)
            
            guard let decrypted else {
                errorMessage = "Decryption failed"
                isFetching = false
                return
            }
            
            ClipboardService.shared.copyToClipboard(data: decrypted)
            isFetching = false
        }
    }
}

// MARK: - Fetch
extension HomeViewModel {
    func fetchEntries() async {
        let (res, err) = await entryService.fetchEntries(deviceID: deviceID)
        try? await Task.sleep(nanoseconds: UInt64(2 * Double(NSEC_PER_SEC)))

        await MainActor.run {
            errorMessage = err
            guard let res else {
                isFetching = false
                return
            }
            
            entries = res
            isFetching = false
        }
    }
}
