//
//  Error.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-08.
//

import Foundation

// MARK: - ErrorResponse
struct ErrorResponse: Codable {
    let error: String
}

// MARK: - APIError
enum APIError: Error {
    case serverError(message: String, statusCode: Int)
    case unknownError
}
