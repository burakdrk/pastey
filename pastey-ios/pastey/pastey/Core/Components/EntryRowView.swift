//
//  EntryRowView.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-11.
//

import SwiftUI

struct EntryRowView: View {
    @State private var isExpanded = false
    @State private var rotationAngle: Double = 0

    let entry: Entry
    
    var body: some View {
        
        VStack {
            HStack {
                Text(isExpanded ? "From \(entry.fromDeviceName)" : entry.encryptedData)
                    .lineLimit(1)
                    .padding(.trailing, 15)
                    .bold(isExpanded ? true : false)
                
                Spacer()
                
                VStack {
                    Image(systemName: "chevron.right")
                        .font(.title2)
                        .rotationEffect(Angle(degrees: rotationAngle))
                        .bold()
                }
                .foregroundColor(Color.theme.accent)
            }
            .padding(.vertical, 15)
            .padding(.horizontal, 5)
            .onTapGesture {
                isExpanded.toggle()
                withAnimation(.easeInOut(duration: 0.3)) {
                    rotationAngle += 90
                }
                rotationAngle = isExpanded ? 90 : 0
            }
            
            if isExpanded {
                Text(entry.encryptedData)
                
                Button("Copy") {
                    print("copy")
                }.foregroundColor(Color.theme.accent)
                    .frame(alignment: .leading)
                    
            }
            
            
        }
        
    }
}

#Preview(traits: .sizeThatFitsLayout) {
    Group {
        EntryRowView(entry: DeveloperPreview.instance.entry)
    }
}
