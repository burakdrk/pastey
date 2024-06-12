//
//  ClipboardService.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-10.
//

import Foundation
import UIKit

class ClipboardService {
    static let shared = ClipboardService()
    
    private init() {}
    
    func copyToClipboard(data: String) {
        UIPasteboard.general.string = data
    }
    
    func pasteFromClipboard() -> String? {
        return UIPasteboard.general.string
    }
}
