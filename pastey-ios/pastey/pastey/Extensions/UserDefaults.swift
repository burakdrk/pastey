//
//  UserDefaults.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-10.
//

import Foundation

extension UserDefaults {
    func setObject<T: Codable>(_ object: T, forKey key: String) {
        let encoder = JSONEncoder()
        if let encoded = try? encoder.encode(object) {
            self.set(encoded, forKey: key)
        }
    }
        
    func getObject<T: Codable>(forKey key: String, as type: T.Type) -> T? {
        if let data = self.data(forKey: key) {
            let decoder = JSONDecoder()
            return try? decoder.decode(T.self, from: data)
        }
        return nil
    }
}
