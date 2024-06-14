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
    
    func fetchUser() async -> (User?, String?) {
        let url = URL(string: "\(API_URL)/users")!
        
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
    
    func logout(refreshToken: String) async -> String? {
        let url = URL(string: "\(API_URL)/users/logout")!
        let request = LogoutRequest(refreshToken: refreshToken)
        
        do {
            _ = try await APIClient.postWithoutResponse(url: url, request: request)
            return nil
        } catch APIError.serverError(message: let message, statusCode: let statusCode) {
            switch statusCode {
            case 401:
                return "Unauthorized"
            default:
                return message
            }
        } catch {
            print("DEBUG: UserService::logout: \(error.localizedDescription)")
            return "An unexpected error has occured. Please try again later."
        }
    }
    
    func createUser(email: String, password: String) async -> (User?, String?) {
        let url = URL(string: "\(API_URL)/users")!
        let request = CreateUserRequest(email: email, password: password)
        
        do {
            return (try await APIClient.post(url: url, request: request), nil)
        } catch APIError.serverError(message: let message, statusCode: let statusCode) {
            switch statusCode {
            case 409:
                return (nil, "User with this email already exists")
            default:
                return (nil, message)
            }
        } catch {
            print("DEBUG: UserService::createUser: \(error.localizedDescription)")
            return (nil, "An unexpected error has occured. Please try again later.")
        }
    }
}

