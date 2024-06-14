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
    
    func fetchDevices() async -> (ListDevicesResponse?, String?) {
        let url = URL(string: "\(API_URL)/devices")!
        
        do {
            return (try await APIClient.fetch(url: url), nil)
        } catch APIError.serverError(message: let message, statusCode: let statusCode) {
            switch statusCode {
            case 401:
                return (nil, "Unauthorized")
            default:
                return (nil, message)
            }
        } catch {
            print("DEBUG: DeviceService::fetchDevices: \(error.localizedDescription)")
            return (nil, "An unexpected error has occured. Please try again later.")
        }
    }
    
    func deleteDevice(deviceID: Int) async -> String? {
        let url = URL(string: "\(API_URL)/devices/\(deviceID)")!
        
        do {
            try await APIClient.delete(url: url)
            return nil
        } catch APIError.serverError(message: let message, statusCode: let statusCode) {
            switch statusCode {
            case 401:
                return "Unauthorized"
            default:
                return message
            }
        } catch {
            print("DEBUG: DeviceService::deleteDevice: \(error.localizedDescription)")
            return "An unexpected error has occured. Please try again later."
        }
    }
}
