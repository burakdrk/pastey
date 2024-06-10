//
//  UserService.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-07.
//

import Foundation

class UserService {
    func login(email: String, password: String, deviceID: Int? = nil) async -> (LoginResponse?, String?) {
        let url = URL(string: "\(API_URL)/users/login")!
        let request = LoginRequest(email: email, password: password, deviceID: deviceID)
        
        do {
            return (try await APIClient.post(url: url, request: request), nil)
        } catch APIError.serverError(message: let message, statusCode: let statusCode) {
            switch statusCode {
            case 404:
                return (nil, "User not found")
            case 401:
                return (nil, "Invalid credentials")
            default:
                return (nil, message)
            }
        } catch {
            print("DEBUG: UserService::login: \(error.localizedDescription)")
            return (nil, "An unexpected error has occured. Please try again later.")
        }
    }
    
    func fetchUser(userID: Int = 0) async -> (User?, String?) {
        let url = URL(string: "\(API_URL)/users/\(userID != 0 ? String(userID) : "me")")!
        
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
            print("DEBUG: UserService::fetchUser: \(error.localizedDescription)")
            return (nil, "An unexpected error has occured. Please try again later.")
        }
    }
}

