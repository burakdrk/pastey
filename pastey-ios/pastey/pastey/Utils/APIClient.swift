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
        let keychain = SimpleKeychain()
        let accessToken: String?
        
        do {
            accessToken = try keychain.string(forKey: "access_token")
        } catch {
            accessToken = nil
        }
        
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
        let keychain = SimpleKeychain()
        let accessToken: String?
        
        do {
            accessToken = try keychain.string(forKey: "access_token")
        } catch {
            accessToken = nil
        }
        
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
}
