//
//  EntryRowView.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-11.
//

import SwiftUI
import AlertKit

struct EntryRowView: View {
    @State private var isExpanded = false
    
    let entry: Entry
    
    var decryptedData: String {
        let decrypt = try? EncryptionService.shared.decryptMessage(encryptedMessage: entry.encryptedData)
        
        return decrypt ?? "Decryption failed"
    }
    
    var body: some View {
        DisclosureGroup(isExpanded: $isExpanded) {
            VStack {
                Text(decryptedData)
                    .padding(.bottom, 20)
                    .frame(maxWidth: .infinity, alignment: .leading)
                    .onTapGesture {
                        ClipboardService.shared.copyToClipboard(data: decryptedData)
                        AlertKitAPI.present(
                            title: "Copied to clipboard",
                            icon: .done,
                            style: .iOS16AppleMusic,
                            haptic: .success
                        )
                    }
                
                HStack {
                    if let date = DateFormatter.iso8601Full.date(from: entry.createdAt) {
                        HStack {
                            Text(date, style: .time)
                            Text("•")
                            Text(date, style: .date)
                        }
                    }
                    Spacer()
                    Text(entry.fromDeviceName)
                }
                .font(.caption)
                .foregroundColor(Color.theme.secondaryText)
                .frame(maxWidth: .infinity, alignment: .leading)
            }
            .listRowSeparator(.hidden, edges: .top)
            .listRowInsets(EdgeInsets(top: 0, leading: 0, bottom: 15, trailing: 20))
            .deleteDisabled(true)
        } label : {
            HStack {
                Text(isExpanded ? "Entry" : decryptedData)
                    .lineLimit(1)
                    .padding(.trailing, 15)
                    .bold(isExpanded ? true : false)
                    .frame(maxWidth: .infinity, alignment: .leading)
                
            }
            .padding(.vertical, 13)
        }
        .alignmentGuide(.listRowSeparatorTrailing) { d in
            d[.trailing]
        }
    }
}

#Preview(traits: .sizeThatFitsLayout) {
    Group {
        EntryRowView(entry: DeveloperPreview.instance.entry)
    }
}
