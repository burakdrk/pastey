//
//  APIClient.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-07.
//

import Foundation
import SimpleKeychain

class APIClient {
    static let decoder = JSONDecoder()
    static let encoder = JSONEncoder()
    
    private init() {}

    static func fetch<T: Decodable>(url: URL) async throws -> T {
        let accessToken = await getAccessToken()
        
        var urlRequest = URLRequest(url: url)
        urlRequest.httpMethod = "GET"
        urlRequest.setValue("application/json", forHTTPHeaderField: "Content-Type")
        if let accessToken {
            urlRequest.setValue("Bearer \(accessToken)", forHTTPHeaderField: "Authorization")
        }

        let (data, response) = try await URLSession.shared.data(for: urlRequest)
        
        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError.unknownError
        }
        
        guard (200...299).contains(httpResponse.statusCode) else {
            let ErrorResponse = try decoder.decode(ErrorResponse.self, from: data)
            throw APIError.serverError(message: ErrorResponse.error, statusCode: httpResponse.statusCode)
        }
        
        return try decoder.decode(T.self, from: data)
    }

    static func post<T: Decodable, U: Encodable>(url: URL, request: U, method: String = "POST") async throws -> T {
        let accessToken = await getAccessToken()
        
        var urlRequest = URLRequest(url: url)
        urlRequest.httpMethod = method
        urlRequest.setValue("application/json", forHTTPHeaderField: "Content-Type")
        if let accessToken {
            urlRequest.setValue("Bearer \(accessToken)", forHTTPHeaderField: "Authorization")
        }
        urlRequest.httpBody = try encoder.encode(request)

        let (data, response) = try await URLSession.shared.data(for: urlRequest)
        
        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError.unknownError
        }
        
        guard (200...299).contains(httpResponse.statusCode) else {
            let ErrorResponse = try decoder.decode(ErrorResponse.self, from: data)
            throw APIError.serverError(message: ErrorResponse.error, statusCode: httpResponse.statusCode)
        }
        
        return try decoder.decode(T.self, from: data)
    }
    
    static func postWithoutResponse<T: Encodable>(url: URL, request: T, method: String = "POST") async throws {
        let accessToken = await getAccessToken()
        
        var urlRequest = URLRequest(url: url)
        urlRequest.httpMethod = method
        urlRequest.setValue("application/json", forHTTPHeaderField: "Content-Type")
        if let accessToken {
            urlRequest.setValue("Bearer \(accessToken)", forHTTPHeaderField: "Authorization")
        }
        urlRequest.httpBody = try encoder.encode(request)
        
        let (data, response) = try await URLSession.shared.data(for: urlRequest)
        
        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError.unknownError
        }
        
        guard (200...299).contains(httpResponse.statusCode) else {
            let ErrorResponse = try decoder.decode(ErrorResponse.self, from: data)
            throw APIError.serverError(message: ErrorResponse.error, statusCode: httpResponse.statusCode)
        }
    }
    
    static func delete(url: URL) async throws {
        let accessToken = await getAccessToken()
        
        var urlRequest = URLRequest(url: url)
        urlRequest.httpMethod = "DELETE"
        urlRequest.setValue("application/json", forHTTPHeaderField: "Content-Type")
        if let accessToken {
            urlRequest.setValue("Bearer \(accessToken)", forHTTPHeaderField: "Authorization")
        }
        
        let (data, response) = try await URLSession.shared.data(for: urlRequest)
        
        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError.unknownError
        }
        
        guard (200...299).contains(httpResponse.statusCode) else {
            let ErrorResponse = try decoder.decode(ErrorResponse.self, from: data)
            throw APIError.serverError(message: ErrorResponse.error, statusCode: httpResponse.statusCode)
        }
    }
    
    private static func getAccessToken() async -> String? {
        let keychain = SimpleKeychain()
        let dateFormatter = DateFormatter.iso8601Full

        var accessToken: String
        let accessTokenExpiry: Date?

        do {
            accessToken = try keychain.string(forKey: "access_token")
            accessTokenExpiry = dateFormatter.date(from: try keychain.string(forKey: "access_token_expiry"))
        } catch {
            return nil
        }
                
        if let accessTokenExpiry {
            if accessTokenExpiry > Date.now {
                return accessToken
            }
        } else {
            return nil
        }
        
        // If the access token is expired, refresh it
        let refreshToken: String

        do {
            refreshToken = try keychain.string(forKey: "refresh_token")
        } catch {
            return nil
        }
                
        do {
            let url = URL(string: "\(API_URL)/token/refresh")!
            let request = RefreshTokenRequest(refreshToken: refreshToken)
            var urlRequest = URLRequest(url: url)
            urlRequest.httpMethod = "POST"
            urlRequest.setValue("application/json", forHTTPHeaderField: "Content-Type")
            urlRequest.httpBody = try encoder.encode(request)
            
            let (data, response) = try await URLSession.shared.data(for: urlRequest)
            
            guard let httpResponse = response as? HTTPURLResponse else {
                return nil
            }
            guard (200...299).contains(httpResponse.statusCode) else {
                return nil
            }
            
            let res = try decoder.decode(RefreshTokenResponse.self, from: data)

            try keychain.set(res.accessToken, forKey: "access_token")
            try keychain.set(res.accessTokenExpiresAt, forKey: "access_token_expiry")
            accessToken = res.accessToken
        } catch {
            return nil
        }
        
        return accessToken
    }
}
