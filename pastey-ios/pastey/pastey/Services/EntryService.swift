//
//  EntryService.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-11.
//

import Foundation

class EntryService {
    func fetchEntries(deviceID: Int) async -> (EntryResponse?, String?) {
        let url = URL(string: "\(API_URL)/devices/\(deviceID)/entries")!
        
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
            print("DEBUG: EntryService::fetchEntries: \(error.localizedDescription)")
            return (nil, "An unexpected error has occured. Please try again later.")
        }
    }
    
    func copyEntries(fromDeviceID: Int, copies: [Copy]) async -> String? {
        let url = URL(string: "\(API_URL)/copy")!
        let request = CopyRequest(fromDeviceID: fromDeviceID, copies: copies)
        
        do {
            try await APIClient.postWithoutResponse(url: url, request: request)
            return nil
        } catch APIError.serverError(message: let message, statusCode: let statusCode) {
            switch statusCode {
            case 401:
                return "Unauthorized"
            default:
                return message
            }
        } catch {
            print("DEBUG: EntryService::copyEntries: \(error.localizedDescription)")
            return "An unexpected error has occured. Please try again later."
        }
    }
    
    func deleteEntry(entryID: String) async -> String? {
        let url = URL(string: "\(API_URL)/entries/\(entryID)")!
        
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
            print("DEBUG: EntryService::deleteEntry: \(error.localizedDescription)")
            return "An unexpected error has occured. Please try again later."
        }
    }
}
