//
//  ProfileViewModel.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-12.
//

import Foundation
import SwiftUI
import AlertKit

class ProfileViewModel: ObservableObject {
    @Published var user: User?
    @Published var devices: [Device]
    @Published var errorMessage: String?
    @Published var isFetching = false
    
    private let userService: UserService
    private let deviceService: DeviceService
    private let dateFormatter: DateFormatter
    
    init() {
        self.userService = UserService()
        self.deviceService = DeviceService()
        self.dateFormatter = DateFormatter.iso8601Full
        self.user = UserDefaults.standard.getObject(forKey: "user", as: User.self)
        self.devices = []
        
        isFetching = true
        
        Task { await fetchDevices() }
    }
    
    private func setErrorMessage(_ message: String?) {
        self.errorMessage = message
        
        guard let errorMessage else {
            return
        }
        
        AlertKitAPI.present(
            title: errorMessage,
            icon: .error,
            style: .iOS16AppleMusic,
            haptic: .error
        )
    }
    
    func deleteDevice(at offsets: IndexSet) {
        isFetching = true
        
        Task { @MainActor in
            let device = devices[offsets.first!]
            let err = await deviceService.deleteDevice(deviceID: device.id)
            
            setErrorMessage(err)
            if err != nil {
                isFetching = false
                return
            }
            
            self.devices.remove(atOffsets: offsets)
            isFetching = false
            AlertKitAPI.present(
                title: "Deleted \(device.deviceName)",
                icon: .done,
                style: .iOS16AppleMusic,
                haptic: .success
            )
        }
    }
    
    func fetchDevices() async {
        let (devices, err) = await deviceService.fetchDevices()
        
        await MainActor.run {
            setErrorMessage(err)
            guard let devices else {
                isFetching = false
                return
            }
            
            self.devices = devices
            isFetching = false
        }
    }
}
