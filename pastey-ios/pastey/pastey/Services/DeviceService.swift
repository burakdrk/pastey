//
//  DeviceService.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-10.
//

import Foundation
import UIKit

class DeviceService {
    func createDevice() async -> (CreateDeviceResponse?, String?) {
        let url = URL(string: "\(API_URL)/devices")!
        let publicKey: String
        
        do {
            let keyPair = try EncryptionService.shared.generateKeyPair()
            publicKey = keyPair.publicKey
        } catch {
            return (nil, "An unexpected error has occured. Please try again later.")
        }
        
        let request = await CreateDeviceRequest(deviceName: UIDevice.current.name, publicKey: publicKey)
        
        do {
            return (try await APIClient.post(url: url, request: request), nil)
        } catch APIError.serverError(message: let message, statusCode: let statusCode) {
            switch statusCode {
            case 401:
                return (nil, "Unauthorized")
            default:
                return (nil, message)
            }
        } catch {
            print("DEBUG: DeviceService::createDevice: \(error.localizedDescription)")
            return (nil, "An unexpected error has occured. Please try again later.")
        }
    }
}
